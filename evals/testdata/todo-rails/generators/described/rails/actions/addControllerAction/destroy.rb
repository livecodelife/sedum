  def destroy
    # One statement. Loading the row first only to discard it makes the delete
    # two interactions, and the row is not needed for anything.
    raise ActiveRecord::RecordNotFound if {{resource|model}}.where(id: params[:id]).delete_all.zero?

    head :no_content
  end
