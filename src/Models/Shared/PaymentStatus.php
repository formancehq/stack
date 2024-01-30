<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


enum PaymentStatus: string
{
    case Pending = 'PENDING';
    case Succeeded = 'SUCCEEDED';
    case Cancelled = 'CANCELLED';
    case Failed = 'FAILED';
    case Expired = 'EXPIRED';
    case Refunded = 'REFUNDED';
    case RefundedFailure = 'REFUNDED_FAILURE';
    case Dispute = 'DISPUTE';
    case DisputeWon = 'DISPUTE_WON';
    case DisputeLost = 'DISPUTE_LOST';
    case Other = 'OTHER';
}
