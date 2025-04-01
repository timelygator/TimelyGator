import React from 'react'
import StatCard from './StatCard'

describe('<StatCard />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<StatCard />)
  })
})