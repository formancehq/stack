<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


enum OrchestrationPaymentType: string
{
    case PayIn = 'PAY-IN';
    case Payout = 'PAYOUT';
    case Transfer = 'TRANSFER';
    case Other = 'OTHER';
}
