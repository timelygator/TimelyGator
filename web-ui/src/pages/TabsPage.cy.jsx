import React from 'react'
import TabsPage from './TabsPage'

describe('<TabsPage />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<TabsPage />)
  })
})