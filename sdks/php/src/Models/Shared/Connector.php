<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


/** The name of the connector. */
enum Connector: string
{
    case STRIPE = 'STRIPE';
    case DUMMY_PAY = 'DUMMY-PAY';
    case WISE = 'WISE';
    case MODULR = 'MODULR';
    case CURRENCY_CLOUD = 'CURRENCY-CLOUD';
    case BANKING_CIRCLE = 'BANKING-CIRCLE';
    case MANGOPAY = 'MANGOPAY';
}
