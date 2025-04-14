import React from 'react'
import OverviewCards from './OverviewCards'

describe('<OverviewCards />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<OverviewCards />)
  })
})