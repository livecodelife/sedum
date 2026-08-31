require "test_helper"

class TodosControllerTest < ActionDispatch::IntegrationTest
  # sedum:anchor:tests
  # sedum:addControllerTest:index {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"index"}}
  test "index returns 200" do
    get "/todos"

    assert_response :success
  end
  # /sedum:addControllerTest:index
  # sedum:addControllerTest:show {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"show"}}
  test "show returns 404 for a todo that does not exist" do
    get "/todos/0"

    assert_response :not_found
  end
  # /sedum:addControllerTest:show
  # sedum:addControllerTest:create {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"create"}}
  test "create returns the todo it stored" do
    post "/todos", params: { title: "write it down" }, as: :json

    assert_response :created
    assert_equal "write it down", response.parsed_body["title"]
  end
  # /sedum:addControllerTest:create
  # sedum:addControllerTest:update {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"update"}}
  test "update applies what the body carried" do
    created = Todo.create!(title: "before")
    put "/todos/#{created.id}", params: { title: "after", completed: true }, as: :json

    assert_response :success
    assert_equal "after", response.parsed_body["title"]
  end
  # /sedum:addControllerTest:update
  # sedum:addControllerTest:destroy {"tier":"owned","record":"eval-todo-rails","kwargs":{"resource":"todo","verb":"destroy"}}
  test "destroy returns 404 for a todo that does not exist" do
    delete "/todos/0"

    assert_response :not_found
  end
  # /sedum:addControllerTest:destroy
end
