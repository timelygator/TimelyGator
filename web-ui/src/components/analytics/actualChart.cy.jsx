import React from 'react'
import actualChart from './actualChart'

describe('<actualChart />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<actualChart />)
  })
})