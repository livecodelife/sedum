def create
  {{collection|instantize}} = {{collection|constantize}}.create!(params)
  render json: {{collection|instantize}}, status: :created
end
