Feature: Showing how cucumboa works

  Scenario: Get an existing Pet
    When the 'getPetById' operation is called with path params:
      | petId | 1234 |
    Then the response status will be '200'
    And the content will have values:
      | id | 1234 |

  Scenario: Request a Pet that doesn't exist
    When the 'getPetById' operation is called with path params:
      | petId | 9876 |
    Then the response status will be '404'
