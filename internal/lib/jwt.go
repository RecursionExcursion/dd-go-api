package lib

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
	JWT Registered Claim Names:
	Claim	Meaning													Handled Automatically?
	iss		Issuer — who issued the token							 Yes
	sub		Subject — whom the token is about (usually a user ID)	 No (but useful for your app logic)
	aud		Audience — intended recipients							 Yes
	exp		Expiration time — when it expires						 Yes
	nbf		Not before don’t accept before this time				 Yes
	iat		Issued at —when it was issued							 No (but useful)
	jti		JWT ID — unique identifier for the token				 No
*/

func CreateJWT(claims map[string]any, exp time.Duration, secret string) (string, error) {

	jwtClaims := jwt.MapClaims{
		"exp": time.Now().Add(exp).Unix(),
	}

	for k, v := range claims {
		if k == "exp" {
			continue
		}
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(token string, secret string) (bool, jwt.MapClaims, error) {
	parsedToken, err := parseToken(token, secret)
	if err != nil {
		return false, nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return parsedToken.Valid, claims, nil
	}

	return parsedToken.Valid, claims, nil
}

func parseToken(t string, secret string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
