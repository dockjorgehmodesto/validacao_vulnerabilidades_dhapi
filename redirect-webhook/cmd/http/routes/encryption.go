package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"github.com/go-chi/chi"
	"net/http"
)

type Encryption struct {
	service encryption.Service
}

func NewEncryptionRouter(service encryption.Service) Encryption {
	return Encryption{service: service}
}

func (e Encryption) Route(r chi.Router) {
	r.Use(middleware.CheckIssuerId)

	r.Route("/key-vaults", func(r chi.Router) {
		r.Get("/", e.findAllWithPagination)
		r.Post("/", e.create)
		r.Route("/{id}", func(idRoutes chi.Router) {
			idRoutes.Use(middleware.ValidateUUID(helpers.PathParamDefault))
			idRoutes.Get("/", e.findById)
			idRoutes.Put("/", e.update)
			idRoutes.Post("/replicate", e.replicate)

			idRoutes.Route("/aes", func(aesRoutes chi.Router) {
				aesRoutes.Post("/encrypt", e.encryptAes)
			})
		})
	})

}

// @Title listKeyVaults
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary List encryption key vaults
// @Description List all the key vaults properties based on filters and pagination
// @Param   page 				query    int      false "Page of the current result set"
// @Param   limit 				query    int      false "Amount of data that should be returned"
// @Param   sort 				query    string   false "Attribute to order the list (e.g. 'status', 'status DESC', 'status ASC', 'status DESC, issuerPublicKey')"
// @Success 200 {object} encryption.KeyVaultResponsePage
// @Success 204 "No Content"
// @Accept json
// @Router /encryption/key-vaults [get]
func (e Encryption) findAllWithPagination(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	var requestQueryParams encryption.KeyVaultRequest
	helpers.BindQuery(r, &requestQueryParams)

	pages := e.service.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort, requestQueryParams.ToMap())

	helpers.JsonResponse(w, encryption.ConvertPagination(pages), http.StatusOK)
}

// @Title createKey
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary Create a key vault
// @Description Creates a new key vault with an AES to encrypt data\
// @Description The input's public key must follow this pattern:\
// @Description -----BEGIN PUBLIC KEY-----\
// @Description {your key here}\
// @Description -----END PUBLIC KEY-----\
// @Description\
// @Description Also, you can't create a key vault if the RSA key has expired\
// @Description Passing the input's validation, the API creates a new randomized HEX array of 32 bytes, to be used for encrypting throughout the API, with the AES algorithm.\
// @Description An issuer can obtain the AES by decrypting the value given in the "/encryption/key-vaults/{id}/aes/encrypt" resource, for more details, check it's documentation.
// @Param   content body  encryption.KeyVaultCreate  true  "Object for persisting the key vault's data"
// @Success 200 {object} encryption.KeyVaultResponse
// @Failure 400 "Bad request"
// @Accept json
// @Router /encryption/key-vaults [post]
func (e Encryption) create(w http.ResponseWriter, r *http.Request) {
	var input encryption.KeyVaultCreate
	helpers.BindJson(r, &input)
	keyVault := e.service.Create(helpers.GetTenant(r), input.ConvertToEntity())
	helpers.JsonResponse(w, encryption.ConvertEntityToResponse(keyVault), http.StatusCreated)
}

// @Title findKeyVaultById
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary Find key vault by ID
// @Description Finds for a key vault with given ID
// @Param 		id 		path	string 		true 	"The identifier for the key vault"
// @Success 200 {object} encryption.KeyVaultResponse
// @Success 204 "No Content"
// @Accept json
// @Router /encryption/key-vaults/{id} [get]
func (e Encryption) findById(w http.ResponseWriter, r *http.Request) {
	vaultId := chi.URLParam(r, helpers.PathParamDefault)
	keyVault := e.service.FindById(helpers.GetTenant(r), vaultId)
	helpers.JsonResponse(w, encryption.ConvertEntityToResponse(keyVault), http.StatusOK)
}

// @Title updateKeyVault
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary Update a key vault
// @Description Updates a key vault public data
// @Param   id      path    string             true "The identifier for the key vault"
// @Param   content body    encryption.KeyVaultUpdate   true "Object for updating a key vault"
// @Success 200 {object} encryption.KeyVaultResponse
// @Failure 204 "Key vault not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /encryption/key-vaults/{id} [put]
func (e Encryption) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	var input encryption.KeyVaultUpdate
	helpers.BindJson(r, &input)
	keyVault := e.service.Update(helpers.GetTenant(r), input.ConvertToEntity(id))
	helpers.JsonResponse(w, encryption.ConvertEntityToResponse(keyVault), http.StatusOK)
}

// @Title replicateKeyVault
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary Replicate a key vault
// @Description Copies a key vault to another issuer.\
// @Description When a group of issuer is controlled by a single entity and this group's message will be received by a single webhook, we replicate their keys so it's easier to manager the encryption flow.\
// @Param   id      path    string             true "The identifier for the key vault"
// @Param   content body    encryption.KeyVaultReplicate  true "Object for replicating a key vault"
// @Success 200 {object} encryption.KeyVaultResponse
// @Failure 204 "Key vault not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /encryption/key-vaults/{id}/replicate [post]
func (e Encryption) replicate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	var input encryption.KeyVaultReplicate
	helpers.BindJson(r, &input)
	keyVault := e.service.Replicate(helpers.GetTenant(r), id, input.IssuerId)
	helpers.JsonResponse(w, encryption.ConvertEntityToResponse(keyVault), http.StatusOK)
}

// @Title encryptAesKey
// @Tags Encryption, Internals
// @Security issuerIdAuthentication
// @Summary Encrypt an AES key
// @Description Fetches a key vault, and encrypts it's main key with the RSA algorithm. When encrypting the main key, it's used the issuer's public as the RSA key.\
// @Description By having the encrypted main key, an issuer can decrypt it using their RSA private key. Then, having the clear key the issuer can decrypt any payload received by their webhooks.\
// @Description The response's "encryptedAes" is the encrypted AES encoded in base64, after getting the decrypted/clear AES Key in base64, make sure to decode it to HEX. The HEX value is the key that should be used for decrypting your webhook's payload.\
// @Description\
// @Description ### Example\
// @Description In this example, the message "ABCDEF" is encrypted with the AES key "AA87SF6B8A7FS6BSF87B78S5AFBADS6A".\
// @Description First we have to get the Hex byte array:\
// @Description [![Decode key to hexadecimal byte array](/static/DecodeKeyToHexByteArray.PNG)](https://cryptii.com/pipes/cU3G3w)\
// @Description\
// @Description With the key in hands, the next step is to decode the base64 message, use the AES decryption algorithm with the key from the last step and a blank IV, then finally we achieve the clear text:\
// @Description [![Decrypt message](/static/DecodeAndDecryptTextWithAes.PNG)](https://cryptii.com/pipes/0SsiJQ)\
// @Description
// @Param   id      path    string             true "The identifier for the key vault"
// @Success 200 {object} encryption.EncryptAesResponse
// @Failure 204 "Key vault not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /encryption/key-vaults/{id}/aes/encrypt [post]
func (e Encryption) encryptAes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	encryptedAes := e.service.EncryptAes(helpers.GetTenant(r), id)
	helpers.JsonResponse(w, encryption.EncryptAesResponse{EncryptedAes: encryptedAes}, http.StatusOK)
}
