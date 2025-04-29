package handler

import (
	"encoding/json"
	//"bytes"
	"time"
	"net/http"
	"net/url"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/square/go-jose/v3"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Checkin struct {
	Id string `json:"id"`
	Date   string `json:"date"`
	Direction string `json:"direction"`
	ReturnTime   string `json:"returnTime"`
	Status string `json:"status"`
}

type AuthorizationCodeFlowReq struct {
    ClientId  string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    Code string `json:"code"`
    RedirectUri string `json:"redirect_uri"`
    GrantType string `json:"grant_type"`
}

const _GOOGLE_APP_CLIENT_ID = "627127621175-td1fqlg7dfkm4bm3ljbi8q9svuoe3f4b.apps.googleusercontent.com"
const _GOOGLE_APP_CLIENT_SECRET = "INPTQn3uLwJxYQ2CRbhhS30w"
const _GOOGLE_APP_AUTHORIZATION_URI = "https://ui-transportecheckin-app.vercel.app/"
const _GOOGLE_APP_GRANT_TYPE = "authorization_code"
const KN_SECURITY_HOLDER = "security_holder"
const KN_AUTHORIZATION = "Authorization"

// Chave para criptografia e descriptografia
var key = []byte("E4JCVNEuWq02sErStzEM1ZvMrzbuUU12")

var (
    clientID     = "627127621175-td1fqlg7dfkm4bm3ljbi8q9svuoe3f4b.apps.googleusercontent.com"
    clientSecret = "INPTQn3uLwJxYQ2CRbhhS30w"
    redirectURI  = "https://ui-transportecheckin-app.vercel.app/"
    authURL      = "https://accounts.google.com/o/oauth2/auth"
    tokenURL     = "https://oauth2.googleapis.com/token"
)

// Decodifica a resposta em uma estrutura Go
type TokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    TokenType   string `json:"token_type"`
}

// Estrutura do Token
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var userInfo struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}


func handleRegisterNewUser(username string, password string) (string, error) {
	hash := hashAndSalt([]byte(password))
	fmt.Print(hash)
	return hash, nil
}

// Função para criptografar um token JWT
func createJWEToken(username string) (string, error) {
	fmt.Println("Tamanho da chave:", len(key))
	// Definir as claims do token
	claims := CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   "sub",
			Audience:  []string{"aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Token expira em 15 minutos
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Criar o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	raw, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	// Criptografar o token JWT encrypter,
	encrypter, err := jose.NewEncrypter(jose.A128CBC_HS256, jose.Recipient{
		Algorithm: jose.DIRECT, Key: key}, nil)
	if err != nil {
		return "", err
	}
	object, err := encrypter.Encrypt([]byte(raw))
	if err != nil {
		return "", err
	}
	key, err := object.CompactSerialize()
	if err != nil {
		return "", err
	}
	return key, nil
}

// Função para descriptografar e verificar o token JWE
func decryptJWEToken(jweToken string) (*CustomClaims, error) {
	object, err := jose.ParseEncrypted(jweToken)
	if err != nil {
		return nil, err
	}
	decrypted, err := object.Decrypt(key)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(string(decrypted), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido ou expirado")
	}
	// Verifica se o token expirou
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expirado")
	}
	return claims, nil
}

// Função para gerar novo access token usando refresh token
func refreshAccessToken(refreshToken string) (string, error) {
	// Validar o refresh token
	claims, err := decryptJWEToken(refreshToken)
	if err != nil {
		return "", err
	}
	// Gerar novo access token válido por 15 minutos
	newAccessToken, err := createJWEToken(claims.Username)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func setCookieHandler(w http.ResponseWriter, r *http.Request, name string, value string, path string) {

	//expiration := time.Now().Add(5 * time.Minute)

	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	// Write a HTTP response as normal.
	//w.Write([]byte("cookie set!"))
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status":"Up"})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","name": "John", "email": "johndoe@gmail.com","userType":"mensalista"})
}

func RegisterCheckin(c *gin.Context) {
	var checkin Checkin
    if err := c.BindJSON(&checkin); err != nil {
        return
    }
    checkin.Id = uuid.New().String()
    checkin.Status = "REGISTERED"
	c.IndentedJSON(http.StatusOK, checkin)
}

func GetAllCheckins(c *gin.Context) {
	checkins := generateFakeCheckins()
    c.JSON(http.StatusOK, checkins)
}

func getCookie(r *http.Request, name string) (string, error) {
	// Read the cookie as normal.
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	// Return the decoded cookie value.
	return string(cookie.Value), nil
}

func GetProfile(c *gin.Context) {
	
	token, _ := getCookie(c.Request, KN_SECURITY_HOLDER)
  
  if token == "" {
    token := c.GetHeader(KN_AUTHORIZATION)
    if token == "" {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
      c.Abort()
      return
    }
  }

  claims, err := decryptJWEToken(token)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
    c.Abort()
    return
  }

  c.JSON(http.StatusOK, gin.H{"username": claims.Username})
}

// Função para trocar o código pelo token
func exchangeCodeForToken(code string) (TokenResponse, error) {
	  var token TokenResponse
    /*data := map[string]string{
        "code":          code,
        "client_id":     clientID,
        "client_secret": clientSecret,
        "redirect_uri":  redirectURI,
        "grant_type":    "authorization_code",
    }*/
    // Cria os dados do formulário
    data := url.Values{}
    data.Set("code", code)
    data.Set("client_id", clientID)
    data.Set("client_secret", clientSecret)
    data.Set("redirect_uri", redirectURI)
    data.Set("grant_type", "authorization_code")

    // Cria a requisição POST
    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, err
    }

    // Define os cabeçalhos adequados
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Executa a requisição
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Lê a resposta
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

		if err := json.Unmarshal(body, &token); err != nil {
	    return nil, err
	  }    

	  return token, nil
}

func RealizarSocialLogin(c *gin.Context) {

	var authorization_code = c.Query("code")

	if (authorization_code == "") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})		
		return
	}

  /*data := AuthorizationCodeFlowReq{
            ClientId:     _GOOGLE_APP_CLIENT_ID,
            ClientSecret: _GOOGLE_APP_CLIENT_SECRET,
            Code:          authorization_code,
            RedirectUri:  _GOOGLE_APP_AUTHORIZATION_URI,
            GrantType:    _GOOGLE_APP_GRANT_TYPE,
  }*/

	token, err := exchangeCodeForToken(authorization_code)
  if err != nil {
    log.Println("Erro ao obter token:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha na autenticação"})
    return
  }
	c.JSON(http.StatusOK, gin.H{"token": token})

  /*} else {
		// Retorna apenas o access_token para o cliente
	  url := "https://www.googleapis.com/oauth2/v3/userinfo"
	  req, err := http.NewRequest("POST", url, nil)
	  if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar a requisição"})
		  c.Abort()
	    return
	  }
	  req.Header.Set("Content-Type", "application/json")
	  req.Header.Set("Authorization", "Bearer" + tokenResponse.AccessToken)

	  client := &http.Client{}
	  resp, err := client.Do(req)
	  if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar a requisição"})
	    c.Abort()
	    return
	  }
		defer resp.Body.Close()

	  // Lê a resposta
	  body, err := ioutil.ReadAll(resp.Body)
	  if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler a resposta"})
	    c.Abort()
	  	return
		}

		if err := json.Unmarshal(body, &userInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar o JSON de resposta"})
			c.Abort()
	    return
	  }

		token, err := createJWEToken(userInfo.Email)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar token"})
			return
		}
		setCookieHandler(c.Writer, c.Request, "security_holder", tokenResponse.AccessToken, "localhost")
		
		decrypted_token, err := decryptJWEToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar token"})
			return			
		}

		c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": userInfo, "expire_in": tokenResponse.ExpiresIn})
		return
  }*/
}

func RealizeLogin(c *gin.Context) {
	var req Login
    if err := c.BindJSON(&req); err != nil {
        return
    }    
	if(req.Password != "" && req.Password == "itau1234"){
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}

// Função para criar checkins fake
func generateFakeCheckins() []Checkin {
    var checkins []Checkin
    currentDate := time.Now()
    date := currentDate.Format("2006-01-02")
    checkins = append(checkins, Checkin{uuid.New().String(), date, "ida", "", "REGISTERED"})
    return checkins
}