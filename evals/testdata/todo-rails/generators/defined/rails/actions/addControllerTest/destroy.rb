  test "destroy returns 404 for a {{resource|record}} that does not exist" do
    delete "/{{resource|table}}/0"

    assert_response :not_found
  end
