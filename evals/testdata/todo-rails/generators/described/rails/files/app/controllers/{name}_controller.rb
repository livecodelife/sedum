class {{name|models}}Controller < ApplicationController
  # A missing record is a 404 once for the whole controller rather than a rescue
  # repeated in every action. Each action finds its own record and none of them
  # handles the absence, which is what keeps an action a single rendered
  # fragment with nothing flowing between them.
  rescue_from ActiveRecord::RecordNotFound do
    render json: { error: "not found" }, status: :not_found
  end

  # sedum:anchor:actions

  private

  # sedum:anchor:private
end
