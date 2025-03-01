import React from 'react'
import TabTable from './TabTable'

describe('<TabTable />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<TabTable />)
  })
})