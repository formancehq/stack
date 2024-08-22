<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
use formance\stack\Utils\SpeakeasyMetadata;
class RunScriptRequest
{
    /**
     *
     * @var Shared\Script $script
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public Shared\Script $script;

    /**
     * Name of the ledger.
     *
     * @var string $ledger
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;

    /**
     * Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker.
     *
     * @var ?bool $preview
     */
    #[SpeakeasyMetadata('queryParam:style=form,explode=true,name=preview')]
    public ?bool $preview = null;

    /**
     * @param  ?Shared\Script  $script
     * @param  ?string  $ledger
     * @param  ?bool  $preview
     */
    public function __construct(?Shared\Script $script = null, ?string $ledger = null, ?bool $preview = null)
    {
        $this->script = $script;
        $this->ledger = $ledger;
        $this->preview = $preview;
    }
}