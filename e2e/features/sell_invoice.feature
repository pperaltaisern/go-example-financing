Feature: sell invoices
  In order to obtain financiation
  As an issuer
  I need to be able to put invoices on sale

  Scenario: Eat 5 out of 12          # features/godogs.feature:6
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining