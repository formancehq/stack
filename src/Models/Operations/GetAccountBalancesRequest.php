<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class GetAccountBalancesRequest
{
    /**
     * The account ID.
     *
     * @var string $accountId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=accountId')]
    public string $accountId;

    /**
     * Filter balances by currency.
     *
     * If not specified, all account's balances will be returned.
     *
     *
     * @var ?string $asset
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=asset')]
    public ?string $asset = null;

    /**
     * Parameter used in pagination requests. Maximum page size is set to 15.
     *
     * Set to the value of next for the next page of results.
     * Set to the value of previous for the previous page of results.
     * No other parameters can be set when this parameter is set.
     *
     *
     * @var ?string $cursor
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=cursor')]
    public ?string $cursor = null;

    /**
     * Filter balances by date.
     *
     * If not specified, all account's balances will be returned.
     *
     *
     * @var ?\DateTime $from
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=from,dateTimeFormat=Y-m-d\TH:i:s.up')]
    public ?\DateTime $from = null;

    /**
     * The maximum number of results to return per page.
     *
     * @var ?int $limit
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=limit')]
    public ?int $limit = null;

    /**
     * The maximum number of results to return per page.
     *
     *
     *
     * @var ?int $pageSize
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pageSize')]
    public ?int $pageSize = null;

    /**
     * Fields used to sort payments (default is date:desc).
     *
     * @var ?array<string> $sort
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=sort')]
    public ?array $sort = null;

    /**
     * Filter balances by date.
     *
     * If not specified, default will be set to now.
     *
     *
     * @var ?\DateTime $to
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=to,dateTimeFormat=Y-m-d\TH:i:s.up')]
    public ?\DateTime $to = null;

    /**
     * @param  ?string  $accountId
     * @param  ?string  $asset
     * @param  ?string  $cursor
     * @param  ?\DateTime  $from
     * @param  ?int  $limit
     * @param  ?int  $pageSize
     * @param  ?array<string>  $sort
     * @param  ?\DateTime  $to
     */
    public function __construct(?string $accountId = null, ?string $asset = null, ?string $cursor = null, ?\DateTime $from = null, ?int $limit = null, ?int $pageSize = null, ?array $sort = null, ?\DateTime $to = null)
    {
        $this->accountId = $accountId;
        $this->asset = $asset;
        $this->cursor = $cursor;
        $this->from = $from;
        $this->limit = $limit;
        $this->pageSize = $pageSize;
        $this->sort = $sort;
        $this->to = $to;
    }
}