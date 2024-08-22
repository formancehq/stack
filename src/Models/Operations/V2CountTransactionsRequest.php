<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class V2CountTransactionsRequest
{
    /**
     * $requestBody
     *
     * @var ?array<string, mixed> $requestBody
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?array $requestBody = null;

    /**
     * Name of the ledger.
     *
     * @var string $ledger
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;

    /**
     *
     * @var ?\DateTime $pit
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pit,dateTimeFormat=Y-m-d\TH:i:s.up')]
    public ?\DateTime $pit = null;

    /**
     * @param  ?string  $ledger
     * @param  ?array<string, mixed>  $requestBody
     * @param  ?\DateTime  $pit
     */
    public function __construct(?string $ledger = null, ?array $requestBody = null, ?\DateTime $pit = null)
    {
        $this->ledger = $ledger;
        $this->requestBody = $requestBody;
        $this->pit = $pit;
    }
}