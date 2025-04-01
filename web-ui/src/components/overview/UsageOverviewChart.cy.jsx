import React from 'react'
import UsageOverviewChart from './UsageOverviewChart'

describe('<UsageOverviewChart />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<UsageOverviewChart />)
  })
})