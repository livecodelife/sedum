class TodosController < ApplicationController
  # A missing record is a 404 once for the whole controller rather than a rescue
  # repeated in every action. Each action finds its own record and none of them
  # handles the absence, which is what keeps an action a single rendered
  # fragment with nothing flowing between them.
  rescue_from ActiveRecord::RecordNotFound do
    render json: { error: "not found" }, status: :not_found
  end

  # sedum:anchor:actions
  # sedum:addControllerAction:index {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"index"}}
  def index
    render json: Todo.order(:id)
  end
  # /sedum:addControllerAction:index
  # sedum:addControllerAction:show {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"show"}}
  def show
    render json: Todo.find(params[:id])
  end
  # /sedum:addControllerAction:show
  # sedum:addControllerAction:create {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"create"}}
  def create
    todo = Todo.new(todo_params)
    if todo.save
      # Re-read rather than render the object still in memory. The write and
      # the read the response is built from are then two interactions, which is
      # what lets a protocol-level contract name them separately - the same
      # shape the Go service in this standard writes for the same reason.
      render json: todo.reload, status: :created
    else
      render json: { errors: todo.errors.full_messages }, status: :unprocessable_entity
    end
  end
  # /sedum:addControllerAction:create
  # sedum:addControllerAction:update {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"update"}}
  def update
    # Only what the caller sent. Strong parameters drop everything unpermitted,
    # so an omitted attribute is simply absent here and keeps its stored value.
    attributes = todo_params

    # Write first, then read what the write produced. Loading the row, mutating
    # it and rendering from memory would spend a read before the write and none
    # after, which describes the same change as one interaction where the
    # contract names two.
    changed = Todo.where(id: params[:id]).update_all(
      attributes.to_h.merge("updated_at" => Time.current)
    )
    raise ActiveRecord::RecordNotFound if changed.zero?

    render json: Todo.find(params[:id])
  end
  # /sedum:addControllerAction:update
  # sedum:addControllerAction:destroy {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"destroy"}}
  def destroy
    # One statement. Loading the row first only to discard it makes the delete
    # two interactions, and the row is not needed for anything.
    raise ActiveRecord::RecordNotFound if Todo.where(id: params[:id]).delete_all.zero?

    head :no_content
  end
  # /sedum:addControllerAction:destroy

  private

  # sedum:anchor:private
  # sedum:definePermittedParams {"tier":"owned","record":"eval-todo-rails","kwargs":{"attributes":"title, completed","resource":"todo"}}
  # Permitting on the root parameters rather than requiring the resource key
  # first means a bare JSON body works whether or not parameter wrapping is
  # configured, and an empty body yields an empty hash instead of raising.
  #
  # The attribute list is split rather than interpolated as Ruby so that the
  # caller writes names and nothing else - no colons, no commas to get right.
  def todo_params
    params.permit(*"title, completed".split(/[,\s]+/).reject(&:empty?).map(&:to_sym))
  end
  # /sedum:definePermittedParams
end
