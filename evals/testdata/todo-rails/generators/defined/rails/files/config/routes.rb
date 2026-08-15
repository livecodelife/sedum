Rails.application.routes.draw do
  # Rails' own routes. They are in the template rather than injected because
  # every service in this standard has them and none of them varies: health is
  # what a load balancer polls, and the PWA routes render app/views/pwa.
  get "up" => "rails/health#show", as: :rails_health_check
  get "service-worker" => "rails/pwa#service_worker", as: :pwa_service_worker
  get "manifest" => "rails/pwa#manifest", as: :pwa_manifest

  # sedum:anchor:routes
end
