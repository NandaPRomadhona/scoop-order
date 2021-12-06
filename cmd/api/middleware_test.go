package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"scoop-order/repository/transactions"
	"testing"
)
//func CreateToken(exp int64) (string, error) {
//	fmt.Println(int(exp))
//	payload := schemas.JWTTokenBody{
//		Organizations: []int{1100276},
//		UserID: 2217630,
//		Exp: int(exp),
//		Roles: []int{1,9},
//		Iss: "SCOOP",
//		Sig: "235296e9cc9dbbb310a273ce5c2d2379117ef735",
//		ExpireTimedelta: 0,
//		UserName: "nanda@gramedia.id",
//		Email: "nanda@gramedia.id",
//		DeviceID: 1,
//	}
//	if err != nil {
//		return "", err
//	}
//
//	jwtToken := jwt2.NewWithClaims(jwt2.SigningMethodHS256, payload)
//	return jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
//}

func TestServerAuthenticate(t *testing.T) {
	type field struct {
		name string
		value string
	}
	testCases := []struct {
		name          string
		setupAuth     field
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: field{name: "auth", value: "JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcmdhbml6YXRpb25zIjpbMTEwMDI3Nl0sInVzZXJfaWQiOjIyMTc2MzAsImV4cCI6MTYzNzgyNjA2NSwicm9sZXMiOlsxLDldLCJpc3MiOiJTQ09PUCIsInNpZyI6IjIzNTI5NmU5Y2M5ZGJiYjMxMGEyNzNjZTVjMmQyMzc5MTE3ZWY3MzUiLCJleHBpcmVfdGltZWRlbHRhIjowLCJ1c2VyX25hbWUiOiJuYW5kYUBncmFtZWRpYS5pZCIsImVtYWlsIjoibmFuZGFAZ3JhbWVkaWEuaWQiLCJkZXZpY2VfaWQiOm51bGx9.amGYbJZmTZ5QkFvuu0vm1Eh0guRS_JrUtF75jmnV0JU"},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	trx := transactions.NewTransaction(testDB, testRedis)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, trx)
			authPath := "/auth"
			server.router.GET(
				authPath,
				server.authenticate(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			//token, _ := CreateToken(time.Now().Add(24 * time.Hour).Unix())
			//fmt.Println(token)
			request.Header.Set("Authorization", tc.setupAuth.value)


			server.router.ServeHTTP(recorder, request)
			fmt.Println(recorder.Result().Status)
			tc.checkResponse(t, recorder)
		})
	}
}

//func TestServer_requireAuthentication(t *testing.T) {
//	type fields struct {
//		transaction transactions.Transaction
//		router      *gin.Engine
//		dbGeo       *geoip2.Reader
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   gin.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			server := &Server{
//				transaction: tt.fields.transaction,
//				router:      tt.fields.router,
//				dbGeo:       tt.fields.dbGeo,
//			}
//			if got := server.requireAuthentication(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("requireAuthentication() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}