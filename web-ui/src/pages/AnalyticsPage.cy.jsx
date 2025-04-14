import React from 'react'
import AnalyticsPage from './AnalyticsPage'

describe('<AnalyticsPage />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<AnalyticsPage />)
  })
})