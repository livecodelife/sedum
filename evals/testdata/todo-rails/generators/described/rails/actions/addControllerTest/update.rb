  test "update rejects an empty body" do
    put "/{{resource|table}}/1", params: {}, as: :json

    assert_response :bad_request
  end
