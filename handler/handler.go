package handler

import (
	"encoding/json"
	"bytes"
	"time"
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/square/go-jose/v3"
	"golang.org/x/crypto/bcrypt"
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

const _GOOGLE_APP_CLIENT_ID = "627127621175-td1fqlg7dfkm4bm3ljbi8q9svuoe3f4b.apps.googleusercontent.com"
const _GOOGLE_APP_CLIENT_SECRET = "INPTQn3uLwJxYQ2CRbhhS30w"
const _GOOGLE_APP_AUTHORIZATION_URI = "https://ui-transportecheckin-app.vercel.app/"
const _GOOGLE_APP_GRANT_TYPE = "authorization_code"

// Chave para criptografia e descriptografia
var key = []byte("E4JCVNEuWq02sErStzEM1ZvMrzbuUU12")

// Decodifica a resposta em uma estrutura Go
var tokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    TokenType   string `json:"token_type"`
}

// Estrutura do Token
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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

func RealizarSocialLogin(c *gin.Context) {

	var authorization_code = c.Query("code")

	if (authorization_code == "") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

  // Monta os campos da requisição diretamente na função
  data := map[string]string{
            "client_id":     _GOOGLE_APP_CLIENT_ID,
            "client_secret": _GOOGLE_APP_CLIENT_SECRET,
            "code":          authorization_code,
            "redirect_uri":  _GOOGLE_APP_AUTHORIZATION_URI,
            "grant_type":    _GOOGLE_APP_GRANT_TYPE,
        }

  // Serializa os dados em JSON
  jsonData, err := json.Marshal(data)
  if err != nil {
 	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao converter os dados para JSON"})
   return
  }

  url := "https://oauth2.googleapis.com/token"
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
  if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar a requisição"})
    return
  }
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar a requisição"})
    return
  }
	defer resp.Body.Close()

  // Lê a resposta
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler a resposta"})
  	return
	}	

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar o JSON de resposta"})
    return
  }

	// Retorna apenas o access_token para o cliente
  c.JSON(http.StatusOK, gin.H{
     "access_token": tokenResponse.AccessToken,
     "expires_in":   tokenResponse.ExpiresIn,
     "token_type":   tokenResponse.TokenType,
  })
	//obtem o client_id, client_secret
	      // Aqui você pode enviar o response.code para o servidor
      // Corpo da requisição
          /*const bodyData = new URLSearchParams();
          bodyData.append("client_id", "627127621175-td1fqlg7dfkm4bm3ljbi8q9svuoe3f4b.apps.googleusercontent.com");
          bodyData.append("client_secret", "INPTQn3uLwJxYQ2CRbhhS30w");
          bodyData.append("code", response.code);
          bodyData.append("redirect_uri", "postmessage");
          bodyData.append("grant_type", "authorization_code");

            // Realizando a requisição POST
            const result = await fetch("https://oauth2.googleapis.com/token", {
              method: "POST",
              headers: {
                "Content-Type": "application/x-www-form-urlencoded",
              },
              body: bodyData.toString(),
            });

            if (result.ok) {
              const data = await result.json();
              console.log("Tokens Recebidos:", data);
            } else {
              console.error("Erro na troca de código:", result.status, await result.text());
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