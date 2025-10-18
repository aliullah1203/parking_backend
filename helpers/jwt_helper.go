package helpers

// func GenerateTokens(userID, role string) (string, string, error) {
// 	accessClaims := &Claims{
// 		UserID: userID,
// 		Role:   role,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
// 		},
// 	}
// 	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	refreshClaims := &Claims{
// 		UserID: userID,
// 		Role:   role,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)), // 7 days
// 		},
// 	}
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return accessToken, refreshToken, nil
// }
