Rails.application.routes.draw do
  # Rails' own health route, in the template rather than injected because every
  # service in this standard has it and it does not vary: it is what a load
  # balancer polls, and `rails new --api` generates it.
  #
  # The PWA routes that used to sit here rendered app/views/pwa, which an
  # API-only application does not have (prov-2026-f5e64f22).
  get "up" => "rails/health#show", as: :rails_health_check

  # sedum:anchor:routes
end
