package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"github.com/go-chi/chi"
	gomodel "gopkg.in/jeevatkm/go-model.v1"
	"net/http"
)

type Tenants struct {
	tenants.Service
}

func NewTenantsRouter(service tenants.Service) Tenants {
	return Tenants{service}
}

func (t Tenants) Route(r chi.Router) {
	r.Get("/", t.findAllWithPagination)
	r.Post("/", t.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middleware.ValidateUUID(helpers.PathParamDefault))
		r.Get("/", t.findById)
		r.Put("/", t.update)
		r.Post("/setup", t.setup)
		r.Post("/replicate", t.replicate)
		r.Delete("/", t.delete)
	})
}

// @Title listTenants
// @Tags Internals
// @Summary List Tenant
// @Description List all the tenants
// @Param   page  query    int      false "Page of the current result set"
// @Param   limit query    int      false "Amount of data that should be returned"
// @Param   sort  query    string   false "Attribute to order the list (e.g. 'brand', 'brand DESC', 'brand ASC', 'brand DESC, status')"
// @Success 200 {array} tenants.TenantResponse
// @Failure 204 "No Content"
// @Accept json
// @Router /internals/tenants [get]
func (t Tenants) findAllWithPagination(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	tenantSlice, _ := t.Service.ListAllByPagination(page.Page, page.Limit, page.Sort)
	_ = tenants.Tenant{}

	resultSlice := make([]tenants.TenantResponse, 0)
	for _, tenant := range tenantSlice {
		resultSlice = append(resultSlice, tenants.ConvertEntityToResponse(tenant))
	}

	helpers.JsonResponse(w, resultSlice, http.StatusOK)
}

// @Title createTenant
// @Tags Internals
// @Summary Create a tenant
// @Description Creates a new tenant allowing a new issuer to use the API. Make sure to set it up after creating it, as it's not this resource's responsibility to run migrations and add the Tenant to the application's context!
// @Param content body    tenants.TenantPersist true "Object for creating a tenant"
// @Success 201 {object} tenants.TenantResponse
// @Failure 400 "Bad Request"
// @Accept json
// @Router /internals/tenants [post]
func (t Tenants) create(w http.ResponseWriter, r *http.Request) {
	var tenantPersist tenants.TenantPersist
	helpers.BindJson(r, &tenantPersist)

	tenant := tenantPersist.ToEntity()
	savedTenant := t.Service.Create(tenant)
	helpers.JsonResponse(w, tenants.ConvertEntityToResponse(savedTenant.(tenants.Tenant)), http.StatusCreated)
}

// @Title replicateTenantById
// @Tags Internals
// @Summary Replicate a tenant by ID
// @Description Finds a tenant by its ID, copies it, associating to a new issuer then saves it.\
// @Description This operation is used for setting up new issuers.
// @Param   id path    string  true "The identifier for the tenant"
// @Param content body    tenants.TenantReplicate true "Object for replicating a tenant"
// @Success 200 {object} tenants.TenantResponse
// @Failure 404 "Not found"
// @Accept json
// @Router /internals/tenants/{id}/replicate [post]
func (t Tenants) replicate(w http.ResponseWriter, r *http.Request) {
	tenantId := chi.URLParam(r, helpers.PathParamDefault)
	var tenantReplicate tenants.TenantReplicate
	helpers.BindJson(r, &tenantReplicate)
	tenant := t.Service.Replicate(tenantId, tenantReplicate.IssuerID)
	helpers.JsonResponse(w, tenants.ConvertEntityToResponse(tenant.(tenants.Tenant)), http.StatusOK)
}

// @Title findTenantById
// @Tags Internals
// @Summary Find tenant by ID
// @Description Finds a tenant by its ID
// @Param   id path    string  true "The identifier for the tenant"
// @Success 200 {object} tenants.TenantResponse
// @Failure 404 "Not found"
// @Accept json
// @Router /internals/tenants/{id} [get]
func (t Tenants) findById(w http.ResponseWriter, r *http.Request) {
	tenantId := chi.URLParam(r, helpers.PathParamDefault)
	tenant := t.Service.FindById(tenantId)
	helpers.JsonResponse(w, tenants.ConvertEntityToResponse(tenant), http.StatusOK)
}

// @Title updateTenant
// @Tags Internals
// @Summary Update a tenant
// @Description Updates a tenant by its ID
// @Param id      path    string    true "The identifier for the tenant"
// @Param content body    tenants.TenantPersist true "Object for updating a tenant"
// @Success 200 {object} tenants.TenantResponse
// @Failure 400 "Bad Request"
// @Failure 404 "Not found"
// @Accept json
// @Router /internals/tenants/{id} [put]
func (t Tenants) update(w http.ResponseWriter, r *http.Request) {
	var tenantUpdate tenants.TenantUpdate
	helpers.BindJson(r, &tenantUpdate)

	tenantId := chi.URLParam(r, helpers.PathParamDefault)
	updateData := tenantUpdate.ConvertToEntity(tenantId)
	tenant := t.Service.FindById(tenantId)
	gomodel.Copy(&tenant, updateData)
	updatedTenant := t.Service.Update(tenantId, &tenant)

	helpers.JsonResponse(w, tenants.ConvertEntityToResponse(updatedTenant), http.StatusOK)
}

// @Title setupTenant
// @Tags Internals
// @Summary Setup migration for tenant
// @Description Starts the migration process for a tenant and adds it to a node's context
// @Param id      path    string    true "The identifier for the tenant"
// @Success 200 {object} tenants.TenantResponse
// @Failure 400 "Bad Request"
// @Failure 404 "Not found"
// @Accept json
// @Router /internals/tenants/{id}/setup [post]
func (t Tenants) setup(w http.ResponseWriter, r *http.Request) {
	tenantId := chi.URLParam(r, helpers.PathParamDefault)
	tenant := t.Service.FindById(tenantId)
	e := t.Service.Startup(tenant, true)
	if e != nil {
		helpers.JsonResponse(w, e, http.StatusBadRequest)
		return
	}

	helpers.JsonResponse(w, tenants.ConvertEntityToResponse(tenant), http.StatusOK)
}

// @Title deleteTenant
// @Tags Internals
// @Summary Delete a tenant
// @Description Deletes a tenant by its ID
// @Param   id path    string true "The identifier for the tenant"
// @Success 200 "No content"
// @Failure 204 "Not found"
// @Failure 400 "Bad Request"
// @Accept json
// @Router /internals/tenants/{id} [delete]
func (t Tenants) delete(w http.ResponseWriter, r *http.Request) {
	tenantId := chi.URLParam(r, helpers.PathParamDefault)
	t.Service.Delete(tenantId)
	w.WriteHeader(http.StatusOK)
}
