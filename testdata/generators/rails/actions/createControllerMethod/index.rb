def index
  {{collection|instantize}} = {{collection|constantize}}.all
  render json: {{collection|instantize}}
end
