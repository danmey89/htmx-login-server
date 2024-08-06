package htmxauthentication

import (
	"crypto/x509/pkix"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", handleIndexRequest).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", handleLoginRequest).Methods("POST", "OPTIONS")
	router.HandleFunc("/servelogin", handleServeLoginRequest).Methods("GET")
	router.HandleFunc("/logout", handleLogoutRequest).Methods("POST", "OPTIONS")
	router.HandleFunc("/signup", handleSignupRequest).Methods("POST", "OPTIONS")

	cert, err := generateTLSCertificate()
	if err != nil {
		log.Fatalf("failed to generate TLS certificate: %w", err)
	}

	server := &http.Server{
		Addr:    ".8080",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("failed to start https server: %w", err)
	}
}

func generateTLSCertificate() (tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Authentication playground server."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(7 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return tls.X509KeyPair(certPEM, keyPEM)
}
