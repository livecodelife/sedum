Rails.application.routes.draw do
  # Rails' own health route, in the template rather than injected because every
  # service in this standard has it and it does not vary: it is what a load
  # balancer polls, and `rails new --api` generates it.
  #
  # The PWA routes that used to sit here rendered app/views/pwa, which an
  # API-only application does not have (prov-2026-f5e64f22).
  get "up" => "rails/health#show", as: :rails_health_check

  # sedum:anchor:routes
  # sedum:addRoute:index {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"index"}}
  get "/todos", to: "todos#index"
  # /sedum:addRoute:index
  # sedum:addRoute:show {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"show"}}
  get "/todos/:id", to: "todos#show"
  # /sedum:addRoute:show
  # sedum:addRoute:create {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"create"}}
  post "/todos", to: "todos#create"
  # /sedum:addRoute:create
  # sedum:addRoute:update {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"update"}}
  put "/todos/:id", to: "todos#update"
  # /sedum:addRoute:update
  # sedum:addRoute:destroy {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"destroy"}}
  delete "/todos/:id", to: "todos#destroy"
  # /sedum:addRoute:destroy
end
