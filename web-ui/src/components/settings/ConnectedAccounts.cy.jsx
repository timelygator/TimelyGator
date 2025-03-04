import React from 'react'
import ConnectedAccounts from './ConnectedAccounts'

describe('<ConnectedAccounts />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<ConnectedAccounts />)
  })
})