Feature: Showing how cucumboa works

  Scenario: Referring to OpenId operations
    # When the 'GetFile' operation is called
    When the 'GetFile' operation is called with path params:
      | id | a2/123 |
    Then the response status will be '200'

  Scenario: Files that don't exist
    When the 'GetFile' operation is called with path params:
      | id | a2/doesntexist |
    Then the response status will be '404'
    And the content will have values:
      | error.code | 1 |



# And the content will be an 'ErrorResponse'
# When the 'GetFile' operation is called with path params:
# And the body contains:
# ```
# {

# }
# ```
# Then the response status will be '200'
# And the error code will be '1234'
# And the content will have values:
#   | id | a2/123 |