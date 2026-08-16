  def show
    render json: {{resource|model}}.find(params[:id])
  end
