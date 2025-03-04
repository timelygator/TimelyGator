import React from 'react'
import DangerZone from './DangerZone'

describe('<DangerZone />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<DangerZone />)
  })
})