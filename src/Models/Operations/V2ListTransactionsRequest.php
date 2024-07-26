<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class V2ListTransactionsRequest
{
    /**
     * $requestBody
     *
     * @var ?array<string, mixed> $requestBody
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?array $requestBody = null;

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

    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=expand')]
    public ?string $expand = null;

    /**
     * Name of the ledger.
     *
     * @var string $ledger
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;

    /**
     * The maximum number of results to return per page.
     *
     *
     *
     * @var ?int $pageSize
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pageSize')]
    public ?int $pageSize = null;

    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pit,dateTimeFormat=Y-m-d\TH:i:s.up')]
    public ?\DateTime $pit = null;

    public function __construct()
    {
        $this->requestBody = null;
        $this->cursor = null;
        $this->expand = null;
        $this->ledger = '';
        $this->pageSize = null;
        $this->pit = null;
    }
}