  def index
    render json: {{resource|model}}.order(:id)
  end
