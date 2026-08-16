class {{name|models}}Controller < ApplicationController
  # Callers are services, not browsers, and every request carries JSON. The
  # token check that protects form posts would reject all of them.
  skip_before_action :verify_authenticity_token

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
