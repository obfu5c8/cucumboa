Feature: Custom scenario setup

  Scenario: Proving that CRUD works
    Given pet '1234' exists with name 'freddie'

    When the 'getPetById' operation is called for pet '1234'
    Then the response status will be '200'
    And the pet will be called 'freddie'

    When the 'deletePet' operation is called for pet '1234'
    Then the response status will be '204'

    When the 'getPetById' operation is called for pet '1234'
    Then the response status will be '404'

