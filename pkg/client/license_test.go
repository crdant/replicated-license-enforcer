package client

import (
    "time"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFieldNotInLicense(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    field, _ := c.GetLicenseField("max_member_count")
    if field != nil {
      t.Fatalf("Expected license field not to be in license")
    }

}

func TestValidExpiresAt(t *testing.T) {
    mockExpiresAtField := `{
        "name": "expires_at",
        "title": "Expiration",
        "description": "License Expiration",
        "value": "2025-06-30T04:00:00Z",
        "valueType": "String",
        "signature": {
            "v1": "UaixeEq1y4C8bVy5xa3dAmGrNS0IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ=="
        }
    }` 

    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockExpiresAtField))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    expiresAt, err := c.GetLicenseField("expires_at")
    if expiresAt == nil {
      t.Fatalf("Expected license field to be returned")
    }

    if err != nil {
        t.Fatalf("Did not expect an error, got %v", err)
    }

    if expiresAt.Name != "expires_at" {
        t.Errorf("Expected field name 'expires_at', got '%s'", expiresAt.Name)
    }
    if expiresAt.Value != "2025-06-30T04:00:00Z" {
        t.Errorf("Expected expiration date '2025-06-30T04:00:00Z', got '%s'", expiresAt.Value)
    }
    if expiresAt.Signature.V1 != "UaixeEq1y4C8bVy5xa3dAmGrNS0IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ==" {
        t.Errorf("Expected signature 'UaixeEq1y4C8bVy5xa3dAmGrNS0IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ==', got '%s'", expiresAt.Signature.V1)
    }
}

func TestTamperedExpiresAt(t *testing.T) {
    mockExpiresAtField := `{
        "name": "expires_at",
        "title": "Expiration",
        "description": "License Expiration",
        "value": "2025-06-30T04:00:00Z",
        "valueType": "String",
        "signature": {
            "v1": "UaixeEq1y4C8bVy5xa3dAmGrNS1IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ=="
        }
    }` 

    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockExpiresAtField))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    _, err := c.GetLicenseField("expires_at")
    if err == nil {
        t.Fatalf("Expected signature validation error")
    }
}

func TestValidExpirationDate(t *testing.T) {
    mockExpiresAtField := `{
        "name": "expires_at",
        "title": "Expiration",
        "description": "License Expiration",
        "value": "2025-06-30T04:00:00Z",
        "valueType": "String",
        "signature": {
            "v1": "UaixeEq1y4C8bVy5xa3dAmGrNS0IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ=="
        }
    }` 

    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockExpiresAtField))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    expirationDate, err := c.GetExpirationDate()
    if err != nil {
        t.Fatalf("Expected expiration date to be returned, got %v", err)
    }

    if expirationDate != time.Date(2025, 6, 30, 4, 0, 0, 0, time.UTC) {
       t.Fatalf("Expected expiration date of June 30, 2025, got %v", expirationDate) 
    }
}

