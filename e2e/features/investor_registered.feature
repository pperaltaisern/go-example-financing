Feature: sync investors
  In order to be synchronized with Onboarding
  As Financing
  I need to be able to store investors created in Onboarding

  Scenario: Eat 5 out of 12
    Given an investor is registered
    When an event is received
    Then there should be a copy of that investor in our storage