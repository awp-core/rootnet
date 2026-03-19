package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/cortexia/rootnet/api/internal/server/handler"
	"github.com/cortexia/rootnet/api/internal/server/ws"
)

// RouterParams groups all dependencies for NewRouter
type RouterParams struct {
	Handler       *handler.Handler
	Hub           *ws.Hub
	RelayHandler  *handler.RelayHandler  // nil if relayer not configured
	VanityHandler *handler.VanityHandler // nil if factory not configured
}

// NewRouter creates and configures the HTTP router, registering all API routes and WebSocket endpoints
func NewRouter(p RouterParams) chi.Router {
	h := p.Handler
	hub := p.Hub
	r := chi.NewRouter()

	// Global middleware
	// IMPORTANT: middleware.RealIP trusts X-Forwarded-For/X-Real-IP headers.
	// This is safe ONLY behind a trusted reverse proxy (nginx) that sets these headers.
	// If the API is directly Internet-facing, rate limiting can be bypassed via header spoofing.
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))

	// API route group
	r.Route("/api", func(r chi.Router) {
		// System
		r.Get("/registry", h.GetRegistry)
		r.Get("/health", h.Health)

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUsers)
			r.Get("/count", h.GetUserCount)
			r.Get("/{address}", h.GetUser)
		})

		// Address lookup
		r.Get("/address/{address}/check", h.CheckAddress)

		// Agent nodes
		r.Route("/agents", func(r chi.Router) {
			r.Get("/by-owner/{owner}", h.GetAgentsByOwner)
			r.Get("/by-owner/{owner}/{agent}", h.GetAgentDetail)
			r.Get("/lookup/{agent}", h.LookupAgent)
			r.Post("/batch-info", h.BatchAgentInfo)
		})

		// Staking
		r.Route("/staking", func(r chi.Router) {
			r.Get("/user/{address}/balance", h.GetBalance)
			r.Get("/user/{address}/positions", h.GetStakePositions)
			r.Get("/user/{address}/allocations", h.GetAllocations)
			r.Get("/user/{address}/pending", h.GetPending)
			r.Get("/user/{address}/frozen", h.GetFrozen)
			r.Get("/agent/{agent}/subnet/{subnetId}", h.GetAgentSubnetStake)
			r.Get("/agent/{agent}/subnets", h.GetAgentSubnets)
			r.Get("/subnet/{subnetId}/total", h.GetSubnetTotalStake)
		})

		// Subnets
		r.Route("/subnets", func(r chi.Router) {
			r.Get("/", h.ListSubnets)
			r.Get("/{subnetId}", h.GetSubnet)
			r.Get("/{subnetId}/skills", h.GetSubnetSkills)
			r.Get("/{subnetId}/earnings", h.GetSubnetEarnings)
			r.Get("/{subnetId}/agents/{agent}", h.GetSubnetAgentInfo)
		})

		// Emission
		r.Route("/emission", func(r chi.Router) {
			r.Get("/current", h.GetCurrentEmission)
			r.Get("/schedule", h.GetEmissionSchedule)
			r.Get("/epochs", h.ListEpochs)
		})

		// Tokens
		r.Route("/tokens", func(r chi.Router) {
			r.Get("/awp", h.GetAWPInfo)
			r.Get("/alpha/{subnetId}", h.GetAlphaInfo)
			r.Get("/alpha/{subnetId}/price", h.GetAlphaPrice)
		})

		// Governance
		r.Route("/governance", func(r chi.Router) {
			r.Get("/proposals", h.ListProposals)
			r.Get("/proposals/{proposalId}", h.GetProposal)
			r.Get("/treasury", h.GetTreasury)
		})
	})

	// WebSocket real-time data push
	r.Get("/ws/live", hub.HandleConnect)

	// Relay endpoints (gasless transaction relay, optional)
	if p.RelayHandler != nil {
		r.Route("/api/relay", func(r chi.Router) {
			r.Post("/register", p.RelayHandler.RelayRegister)
			r.Post("/bind", p.RelayHandler.RelayBind)
			r.Post("/set-reward-recipient", p.RelayHandler.RelaySetRewardRecipient)
			r.Post("/allocate", p.RelayHandler.RelayAllocate)
			r.Post("/deallocate", p.RelayHandler.RelayDeallocate)
			r.Post("/activate-subnet", p.RelayHandler.RelayActivateSubnet)
			r.Post("/register-subnet", p.RelayHandler.RelayRegisterSubnet)
		})
	}

	// Vanity salt management + computation (optional)
	// Vanity salt management (always available if factory configured)
	r.Route("/api/vanity", func(r chi.Router) {
		r.Get("/mining-params", h.GetMiningParams)
		r.Post("/upload-salts", h.UploadSalts)
		r.Get("/salts", h.ListAvailableSalts)
		r.Get("/salts/count", h.CountAvailableSalts)
		if p.VanityHandler != nil {
			r.Post("/compute-salt", p.VanityHandler.ComputeSalt)
		}
	})

	return r
}
