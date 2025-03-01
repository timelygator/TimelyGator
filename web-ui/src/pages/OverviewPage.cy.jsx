import React from 'react'
import OverviewPage from './OverviewPage'

describe('<OverviewPage />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<OverviewPage />)
  })
})