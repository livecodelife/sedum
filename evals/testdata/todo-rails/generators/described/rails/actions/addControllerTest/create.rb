  test "create rejects a body with no usable attributes" do
    post "/{{resource|table}}", params: {}, as: :json

    assert_response :unprocessable_entity
  end
