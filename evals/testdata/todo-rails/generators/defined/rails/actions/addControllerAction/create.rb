  def create
    {{resource|record}} = {{resource|model}}.new({{resource|record}}_params)
    if {{resource|record}}.save
      # Re-read rather than render the object still in memory. The write and
      # the read the response is built from are then two interactions, which is
      # what lets a protocol-level contract name them separately - the same
      # shape the Go service in this standard writes for the same reason.
      render json: {{resource|record}}.reload, status: :created
    else
      render json: { errors: {{resource|record}}.errors.full_messages }, status: :unprocessable_entity
    end
  end
