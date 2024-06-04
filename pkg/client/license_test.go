package client

import (
    "time"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

)

func TestIntegerField(t *testing.T) {
    mockMaxMemberCountField := `{
      "name": "member_count_max",
      "title": "Max Member Count",
      "value": 100,
      "valueType": "Integer",
      "signature": {
        "v1": "yMGjD6CcXwnSpqKbWkTdypp319TDkyZJYtr1SOsMDfGN3FAu0XsK+jPgqvuWQcWeDhI31zhjp3305bSgouxLlYCku398/vYLJ5dlZBlfBmzbWMc7yxKE5lyW+PWu6f9KZpw+0uYnQn47t3/5pMvcpk9SVYKjRkmRGKV5kdkPq0SByjcAZFSfO4hLd+Y2zFJB1rLb1z9xKtjPikrwOC2uGEI7pKhkLNmUcgvSyGAPa11xCYbXDIDF3AEMBD5uEwvI9XdfIWmwUvpH7XMPq1SkenMPsSRat4aBx7x+fqqUA7pU5MncoH0ATfd6sC9fv/vj7ZqIvNdZjA3AL1fn+qBdBQ=="
      }
    }`

    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockMaxMemberCountField))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    field, err := c.GetLicenseField("member_count_max")
    assert.NoError(t, err)
    require.NotNil(t, field)

    assert.Equal(t, "member_count_max", field.Name)
    assert.Equal(t, "Max Member Count", field.Title)
    assert.Equal(t, float64(100), field.Value)
    assert.Equal(t, "Integer", field.ValueType)
}

func TestBooleanField(t *testing.T) {
    mockMaxMemberCountField := `{
      "name": "enable_discourse",
      "title": "Enable Discourse (alpha)",
      "value": true,
      "valueType": "Boolean",
      "signature": {
        "v1": "n0BkylfyJ6TYngkJIDMUqh51nw9hDHkA/HKoNd8uE6ADM3E5hW+HdxRHQJaHRbtoYwdwAiF+IrSGdzHuIy1E7KXFvebmNd/5WIdUrWGHjFnzjO3aAoeMZhZC0hyLiBuD4Wfg21p/pf0y5OJav9rGc0G+gablBQ539Okl2jGdOLBSSZrhqYyweaLCsXHkNn7o1Zagl0B9nW7s25juZmNOybAhj6/Yf/8+Lc9CfvE+WM/Q9nE7aT53v65ogkyjenhhJ1+rt2pmrtwMKAgNQQn0U+jGe1nxVv36prLWO5Ncy0Zb3LWtr9va+pgKP0GBUHvC2xqpR7JWBndkrDd72LjZjw=="
      }
    }`

    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockMaxMemberCountField))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    field, err := c.GetLicenseField("member_count_max")
    assert.NoError(t, err)
    require.NotNil(t, field)

    assert.Equal(t, "enable_discourse", field.Name)
    assert.Equal(t, "Enable Discourse (alpha)", field.Title)
    assert.Equal(t, true, field.Value)
    assert.Equal(t, "Boolean", field.ValueType)
}

func TestMissingField(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    field, _ := c.GetLicenseField("max_member_count")
    assert.Nil(t, field)
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
    assert.NoError(t, err)
    assert.NotNil(t, expiresAt)
    assert.Equal(t, "expires_at", expiresAt.Name)
    assert.Equal(t, "2025-06-30T04:00:00Z", expiresAt.Value)
    assert.Equal(t, "UaixeEq1y4C8bVy5xa3dAmGrNS0IdVAWlbJR+p/gsVv3XyeFhEVrHufJxUSKu7hiO/GewtsP8Bv8Cj5mlOnGye/OG4SVhSxSP6gp8yRDiHT0uFnng6eWDqoai3MI9E/GqiUnSgN5ezhN5SdR11KoXm1oGN+YOoPC12rviR8I4jWv9A5Hxv6RSrQUeTUgemw8KweNcT5zXQdmv6xL24dQnnHN9DhiXFxy4nc6ib6qyR8wI7doU2D/xujQIzIcbA7rE1UkUsXSvdRII4EqSiyfz1UDMjerHj3SvG7XSRPLIgr2sXzuXKBP3CgTVBUlKoZ6sPcMSAlutnxEBlNMWHzpfQ==", expiresAt.Signature.V1)
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
    assert.Error(t, err)
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
    assert.NoError(t, err)
    assert.Equal(t, time.Date(2025, 6, 30, 4, 0, 0, 0, time.UTC), expirationDate) 
}
