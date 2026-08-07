def show
  {{collection|instantize}} = {{collection|constantize}}.find(params[:id])
  render json: {{collection|instantize}}
end
