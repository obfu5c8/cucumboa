Feature: Showing how cucumboa works

  Scenario: Vanilla
    When the 'getPetById' operation is called with path params:
      | petId | 1234 |
    Then the response status will be '200'
    And the content will have values:
      | name | doggie |


  Scenario: With custom DSL
    When the 'getPetById' operation is called for pet '1234'
    Then the response status will be '200'
    And the pet will be called 'doggie'

