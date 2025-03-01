import React from 'react'
import Security from './Security'

describe('<Security />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<Security />)
  })
})