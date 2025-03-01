import React from 'react'
import CategoryDistributionChart from './CategoryDistributionChart'

describe('<CategoryDistributionChart />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<CategoryDistributionChart />)
  })
})