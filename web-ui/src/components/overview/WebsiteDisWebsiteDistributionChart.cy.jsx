import React from 'react'
import WebsiteDistributionChart from './WebsiteDis'

describe('<WebsiteDistributionChart />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<WebsiteDistributionChart />)
  })
})