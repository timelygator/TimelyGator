import React from 'react'
import ToggleSwitch from './ToggleSwitch'

describe('<ToggleSwitch />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<ToggleSwitch />)
  })
})