  test "index returns 200" do
    get "/{{resource|table}}"

    assert_response :success
  end
